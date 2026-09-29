package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
)

//struktur data yang harus sama dengan yang dikirim JSON FE
type RequestData struct {
	Nama string `json: "Nama"`
	Email string `json: "Email"`
}

type ResponseData struct{
	Pesan string `json: "Pesan"` //menggunakan "pesan" sesuai alert(data.pesan) di FE 
}

func main(){
	fs := http.FileServer(http.Dir(".")) //karena file html disimpan di folder yg sama dengan golang
	//bisa juga(http.Dir("./frontend")) kalau file html disimpan dlm folder bernama frontend
	http.Handle("/", fs)

	//menghubungkan endpoint '/buy' dengan fungsi buyHandler
	http.HandleFunc("/buy", buyHandler)

	fmt.Println("Server backend Golang berjalan di http://localhost:8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil{
			log.Fatal(err)
		}
}

func buyHandler(w http.ResponseWriter, r *http.Request){
//pastikan metode yang masuk adalah POST
		if r.Method !=http.MethodPost{
			http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
			return
		}

		var data RequestData
		//membaca data JSON yang dikirim oleh JavaScript(fetch)
		err:= json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, errMail := mail.ParseAddress(data.Email)

		if errMail != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseData{
			Pesan: "Gagal! Format email Anda tidak valid(contoh: nama@domain.com).",
			})
			return
		}

		namaTampilan := data.Nama
		if namaTampilan == "" {
			namaTampilan = "(Tanpa Nama)"
		}
	
		//cetak log ke terminal dulu
		fmt.Printf("Menerima pendaftaran! Pembeli: %s | Email: %s\n", namaTampilan, data.Email)

		// (Nanti di baris ini Anda akan menulis kode integrasi ke Payment Gateway seperti Midtrans/Xendit)
   		// (Serta kode untuk simpan data pembeli ke database)

		if data.Nama != "" {
			fmt.Printf("Email: %s --> Nama diisi. Siapkan PDF dengan WATERMARK atas nama: %s\n", data.Email, data.Nama)
		} else {
			fmt.Printf("Email: %s --> Nama kosong. Siapkan PDF STANDAR tanpa watermark.\n", data.Email)
		}

		w.Header().Set("Content-Type", "application/json")
		response := ResponseData{
			Pesan: "Permintaan pembayaran berhasil diproses!",
		}
		json.NewEncoder(w).Encode(response)
	}

	//proses pembayaran novel disini
		//misalnya: integrasi ke payment gateway, simpan ke database, dll)
		//fmt.Printf("Menerima pendaftaran email pembelian: %s\n", data.Email)