package customer

import (
	"context"
	"eltropy/helper"
	"eltropy/model"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddCustomer - Employee can add customer to the bank
func AddCustomer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// Extract token
	ad, err := helper.ExtractToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	// Fetch auth from redis
	username, err := helper.FetchAuth(ad)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	values := r.URL.Query()
	empID := values.Get("empId")
	if username == empID {
		// Add Customer
		var cust model.Customer
		if err := json.NewDecoder(r.Body).Decode(&cust); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		client, _ := helper.GetMongoClient()
		custCollection := client.Database(helper.DB).Collection(helper.CustomerCollection)
		sr := custCollection.FindOne(context.TODO(), bson.M{"custId": cust.CustId})
		if sr.Err() == nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusConflict,
				Msg:  "customer already exist",
			})
		} else {
			_, err := custCollection.InsertOne(context.TODO(), cust)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusCreated,
				Msg:  "Customer added successfully.",
			})
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
	}
}

// DeleteCustomer - Employee can delete customer from bank
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// Extract token
	ad, err := helper.ExtractToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	// Fetch auth from redis
	username, err := helper.FetchAuth(ad)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	values := r.URL.Query()
	empID := values.Get("empId")
	if empID == username {
		// Delete customer
		var id model.Id
		json.NewDecoder(r.Body).Decode(&id)
		client, _ := helper.GetMongoClient()
		collection := client.Database(helper.DB).Collection(helper.CustomerCollection)
		sr := collection.FindOne(context.TODO(), bson.M{"custId": id.Id})
		if sr.Err() == nil {
			collection.DeleteOne(context.TODO(), bson.M{"custId": id.Id})
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusOK,
				Msg:  "customer deleted. id=" + id.Id,
			})
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusNotFound,
				Msg:  "customer not found",
			})
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
	}
}

// UpdateKYC - Update kyc for customer
func UpdateKYC(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// Extract token
	ad, err := helper.ExtractToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	// Fetch auth from redis
	username, err := helper.FetchAuth(ad)
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
		})
		return
	}

	values := r.URL.Query()
	empID := values.Get("empId")
	custID := values.Get("custId")
	if empID == username {
		// Check for customer
		client, _ := helper.GetMongoClient()
		collection := client.Database(helper.DB).Collection(helper.CustomerCollection)
		sr := collection.FindOne(context.TODO(), bson.M{"custId": custID})
		if sr.Err() == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusNotFound,
				Msg:  "Customer not found, id - " + custID,
			})
			return
		}
		// Upload documents
		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("kycDoc")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusBadRequest,
				Msg:  "Error retreiving the file",
			})
			return
		}
		defer file.Close()
		//Check kyc folder exist
		if _, err := os.Stat("." + string(filepath.Separator) + "kyc"); os.IsNotExist(err) {
			os.Mkdir("."+string(filepath.Separator)+"kyc", 0777)
		}
		filename := "kyc/" + custID + "_" + handler.Filename
		// Create file
		dst, err := os.Create(filename)
		defer dst.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusInternalServerError,
				Msg:  "error in creating file.",
			})
			return
		}
		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusInternalServerError,
				Msg:  "error in copy file",
			})
			return
		}
		var customer model.Customer
		sr.Decode(&customer)
		result, err := collection.UpdateOne(
			context.TODO(),
			bson.M{"custId": custID},
			bson.D{
				{"$set", bson.D{{"kyc", model.Kyc{
					IsDone: true,
					KycDoc: filename,
				}}}},
			},
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusInternalServerError,
				Msg:  "Error in saving doc.",
			})
			return
		}
		if result.ModifiedCount >= 1 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(model.Response{
				Code: http.StatusOK,
				Msg:  "kyc doc uploaded successfully.",
			})
			return
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusConflict,
			Msg:  "kyc not updated",
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "unauthorized",
		})
		return
	}
}
