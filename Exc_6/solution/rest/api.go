package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"ordersystem/httptools"
	"ordersystem/model"
	"ordersystem/repository"
	"ordersystem/storage"
	"strings"

	"github.com/go-chi/render"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// GetMenu 			godoc
// @tags 			Menu
// @Description 	Returns the menu of all drinks
// @Produce  		json
// @Success 		200 {array} model.Drink
// @Failure     	500
// @Router 			/api/menu [get]
func GetMenu(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allDrinks, err := db.GetDrinks() // fetch all drinks from the database so we can return the current menu
		if err != nil {
			slog.Error("Unable to load drinks", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to load drinks")
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, allDrinks)
	}
}

// GetOrders		godoc
// @tags 			Order
// @Description 	Returns all orders
// @Produce  		json
// @Success 		200 {array} model.Order
// @Failure     	500
// @Router 			/api/order/all [get]
func GetOrders(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allOrders, err := db.GetOrders() // load all existing orders from the database
		if err != nil {
			slog.Error("Unable to load orders", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to load order")
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, allOrders)
	}
}

// GetOrdersTotal		godoc
// @tags 				Order
// @Description 		Gets totalled orders
// @Produce  			json
// @Success 			200
// @Failure     		500
// @Router 				/api/order/totalled [get]
func GetOrdersTotal(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		totalledOrders, err := db.GetTotalledOrders() // get aggregated totals per drink from the database
		if err != nil {
			slog.Error("Unable to load order totals", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to load order totals")
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, totalledOrders)
	}
}

func GetReceiptFile(db *repository.DatabaseHandler, s3 *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		uintId, err := httptools.ParseUintUrlParam("orderId", r) // extract numeric orderId from URL path to determine which receipt we need to fetch
		if err != nil {
			slog.Error("No order id set") // log this error for debugging when a user does not provide a valid order id in URL
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "No order id set")
			return
		}

		order, err := db.GetOrder(uintId) // retrieve the target order from DB to ensure it exists before looking for its receipt in S3
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, "This order does not exist")
				return
			}
			slog.Error("Unable to load order", slog.String("error", err.Error())) // log database lookup failures for debugging visibility
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to load order")
			return
		}

		ctx := context.Background() // create a context for MinIO communication to manage request lifetime and cancellation

		object, err := s3.GetObject( // attempt to fetch the markdown receipt file from MinIO; retrieves it as a stream from the correct S3 bucket
			ctx,
			storage.OrdersBucket,     // bucket dedicated to storing order receipts to keep storage organized and structured
			order.GetFilename(),      // compute the exact filename (e.g., order_15.md) using the order ID
			minio.GetObjectOptions{}, // default read options since we are simply downloading the object
		)
		if err != nil {
			slog.Error("Unable to load receipt from s3", slog.String("error", err.Error())) // log S3 retrieval failures for diagnostics
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to load receipt")
			return
		}
		defer object.Close() // always close the stream to release system and network resources

		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")                                       // explicitly inform browser that content is markdown for correct rendering or download
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, order.GetFilename())) // suggest appropriate filename to user when downloading

		if _, err := io.Copy(w, object); err != nil { // stream the S3 object directly to the HTTP response without loading everything into memory
			slog.Error("Streaming error", slog.String("error", err.Error()))
		}
	}
}

func PostOrder(db *repository.DatabaseHandler, s3 *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var order model.Order

		payload, err := io.ReadAll(r.Body) // read JSON body from HTTP request; this contains user-submitted order data
		if err != nil {
			slog.Error("Unable to read request body", slog.String("error", err.Error())) // log body read error for debugging
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "Unable to read body")
			return
		}

		err = json.Unmarshal(payload, &order) // unmarshal JSON payload into Go struct so we can validate and store it
		if err != nil {
			slog.Error("Unable to decode request JSON", slog.String("error", err.Error())) // log JSON decode error for visibility
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "Unable to decode body")
			return
		}

		dbOrder, err := db.AddOrder(&order) // insert the parsed order into the database; assigns its unique ID used for the filename
		if err != nil {
			slog.Error("Unable to add order to db", slog.String("error", err.Error())) // 🔥 RESTORED log line (important for DB failure diagnostics)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to add order to db")
			return
		}

		ctx := context.Background()          // context for MinIO upload operation
		content := dbOrder.ToMarkdown()      // create markdown-formatted receipt text using order fields such as date, drink ID, and quantity
		reader := strings.NewReader(content) // convert markdown string to an io.Reader to stream it into MinIO
		size := int64(len(content))          // compute byte-size of markdown text for S3 upload metadata
		filename := dbOrder.GetFilename()    // generate deterministic filename like order_7.md tied to its DB ID

		_, err = s3.PutObject( // upload the receipt file to MinIO so that it becomes permanently stored and retrievable
			ctx,
			storage.OrdersBucket, // bucket used specifically for storing receipts to avoid mixing unrelated objects
			filename,             // the key (filename) under which the receipt will be saved inside the S3 bucket
			reader,               // markdown content stream
			size,                 // exact file size needed by MinIO/S3 during upload
			minio.PutObjectOptions{ContentType: "text/markdown"}, // set the appropriate content type for proper retrieval later
		)
		if err != nil {
			slog.Error("Unable to store receipt in s3", slog.String("error", err.Error())) // 🔥 important: logs MinIO write failure
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "Unable to store receipt")
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, "ok") // acknowledge successful DB save AND receipt upload
	}
}
