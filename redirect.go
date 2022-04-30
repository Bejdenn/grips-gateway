package gateway

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
)

const gripsBaseURL = "https://elearning.uni-regensburg.de/"

type Course struct {
	Id    string   `firebase:"id"`
	Name  string   `firebase:"name"`
	Hints []string `firebase:"hints"`
}

type Handler struct {
}

type Database struct {
	Instance         *firestore.Client
	TargetCollection string
}

func (h *Handler) Handle(db *Database, rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error while parsing request form: %v", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		h := r.FormValue("hint")
		log.Printf("Registered GRIPS gateway request with course hint: %s\n", h)

		iter := db.Instance.Collection(db.TargetCollection).Where("hints", "array-contains", h).Limit(1).Documents(context.TODO())
		defer iter.Stop()
		doc, err := iter.Next()

		if err == iterator.Done {
			log.Printf("No course exists with the a hint '%s'. "+
				"Either there is a typo or the course has to be extended by this particular hint.\n", h)

			http.Redirect(rw, r, gripsBaseURL, http.StatusMovedPermanently)
			return
		} else if err != nil {
			log.Print(err)
		} else {
			var c Course
			err = doc.DataTo(&c)
			if err != nil {
				log.Printf("Could not cast Firestore object to struct: %v", err)
				return
			}

			dest := gripsBaseURL + "course/view.php?id=" + c.Id
			log.Printf("Found course '%s' with ID %s. Redirecting now to %s\n", c.Name, c.Id, dest)
			http.Redirect(rw, r, dest, http.StatusMovedPermanently)
		}

	default:
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func Handle(rw http.ResponseWriter, r *http.Request) {
	projectID := "grips-gateway"

	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	db := &Database{Instance: client, TargetCollection: "courses"}
	h := Handler{}
	h.Handle(db, rw, r)
}
