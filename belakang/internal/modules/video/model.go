package video

import "time"

type Video struct {
	ID           string    `db:"id" json:"id"`
	Nama         string    `db:"nama" json:"nama"`
	Kegiatan     string    `db:"kegiatan" json:"kegiatan"`
	Video        string    `db:"video" json:"video"`
	Thumbnail    string    `json:"thumbnail"`
	YoutubeURL   string    `json:"youtube_url"`
	Waktu        time.Time `db:"waktu" json:"waktu"`
	IDPengguna   string    `db:"id_pengguna" json:"id_pengguna"`
	NamaPengguna string    `db:"nama_pengguna" json:"nama_pengguna"`
}
