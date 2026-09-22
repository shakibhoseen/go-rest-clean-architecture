package handlers

import (
	"context"
	"crud/utils"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

// Context Key টাইপ যাতে অন্য কোনো কি-এর সাথে সংঘর্ষ না হয়
type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorJSON(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// হেডার ফরম্যাট: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorJSON(w, http.StatusUnauthorized, "Invalid token format, expected 'Bearer <token>'")
			return
		}

		tokenString := parts[1]

		// এই যে এখানে ValidateJWT ব্যবহার হলো!
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.ErrorJSON(w, http.StatusUnauthorized, err.Error())
			return
		}

		// টোকেনের পে-লোড থেকে user_id নেওয়া
		userIDFloat, ok := claims["user_id"].(float64) // JSON সংখ্যাক্রম float64 হিসেবে পার্স হয়
		if !ok {
			utils.ErrorJSON(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// রিকোয়েস্টের Context-এ user_id ইনজেক্ট করে পরবর্তী হ্যান্ডলারে পাঠানো
		ctx := context.WithValue(r.Context(), UserIDKey, int(userIDFloat))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// সাধারণ func(w, r)-কে সরাসরি প্রোটেক্ট করার র্যাপার
func Protected(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AuthMiddleware(handlerFunc).ServeHTTP(w, r)
	}
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// সব ডোমেইন থেকে অ্যাক্সেস অ্যালাউ করা
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// অনুমোদিত মেথডগুলো
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// ক্লায়েন্ট থেকে যে হেডারগুলো পাঠানো অনুমোদিত (Authorization মাস্ট)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// ব্রাউজারের পাঠানো Preflight (OPTIONS) রিকোয়েস্ট হলে সাথে সাথে 200 OK দিয়ে রিটার্ন করা
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Saluadan ti server tapno saan nga ag-crash no adda runtime panic
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 1. I-print ti kompleto a stack trace iti terminal tapno makita ti developers
				log.Printf("[PANIC RECOVERED] %v\nStack Trace:\n%s", err, debug.Stack())

				// 2. Isubli ti nadalus a JSON response iti client (awan ti sensitive details)
				utils.ErrorJSON(w, http.StatusInternalServerError, "Internal server error. Please try again later.")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
