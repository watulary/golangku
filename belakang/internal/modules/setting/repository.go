package setting

import "sekolah/internal/database"

// =====================================
// Setting Halaman Publik
// =====================================

func GetSetting() (*Setting, error) {

	var setting Setting

	err := database.DB.Get(
		&setting,
		`SELECT
			id_sekolah,
			nama_sekolah,
			alamat_sekolah,
			telp_sekolah,
			email_sekolah,
			logo_sekolah,
			facebook,
			instagram,
			youtube,
			twitter,
			web_sekolah,
			peta_sekolah,
			visi_sekolah,
			deskripsi_sekolah,
			kelurahan_sekolah,
			kecamatan_sekolah,
			kabupaten_sekolah,
			provinsi_sekolah,
			kodepos_sekolah,
			profil_sekolah,
			heroarea1,
			heroarea2,
			heroarea3,
			heroarea4,
			footer_sekolah
		FROM setting
		LIMIT 1`,
	)

	if err != nil {
		return nil, err
	}

	return &setting, nil
}

// =====================================
// Dashboard Setting
// =====================================

func GetDashboardSetting() (*Setting, error) {

	var setting Setting

	err := database.DB.Get(
		&setting,
		`
		SELECT
			s.id_sekolah,
			s.nama_sekolah,
			s.alamat_sekolah,
			s.npsn,
			s.nsss,
			s.nsp,
			s.kelurahan_sekolah,
			s.kecamatan_sekolah,
			s.kabupaten_sekolah,
			s.provinsi_sekolah,
			s.kodepos_sekolah,
			s.telp_sekolah,
			s.email_sekolah,
			s.web_sekolah,
			s.deskripsi_sekolah,
			s.visi_sekolah,
			s.footer_sekolah,
			s.profil_sekolah,
			s.facebook,
			s.twitter,
			s.instagram,
			s.youtube,
			s.peta_sekolah,
			s.logo_sekolah,
			s.heroarea1,
			s.heroarea2,
			s.heroarea3,
			s.heroarea4,
s.id_pengguna AS updated_by,
s.waktu AS updated_at,
p.nama_pengguna AS updated_by_nama

FROM setting s

LEFT JOIN pengguna p
ON p.id_pengguna = s.id_pengguna

		LIMIT 1
		`,
	)

	if err != nil {
		return nil, err
	}

	return &setting, nil
}

// =====================================
// Update Setting
// =====================================

func UpdateSetting(setting *Setting) (bool, error) {

	oldSetting, err := GetDashboardSetting()
	if err != nil {
		return false, err
	}

	// =====================================
	// Tidak Ada Perubahan
	// =====================================

	if oldSetting.NamaSekolah == setting.NamaSekolah &&
		oldSetting.AlamatSekolah == setting.AlamatSekolah &&
		oldSetting.NPSN == setting.NPSN &&
		oldSetting.NSSS == setting.NSSS &&
		oldSetting.NSP == setting.NSP &&
		oldSetting.KelurahanSekolah == setting.KelurahanSekolah &&
		oldSetting.KecamatanSekolah == setting.KecamatanSekolah &&
		oldSetting.KabupatenSekolah == setting.KabupatenSekolah &&
		oldSetting.ProvinsiSekolah == setting.ProvinsiSekolah &&
		oldSetting.KodePosSekolah == setting.KodePosSekolah &&
		oldSetting.TelpSekolah == setting.TelpSekolah &&
		oldSetting.EmailSekolah == setting.EmailSekolah &&
		oldSetting.WebSekolah == setting.WebSekolah &&
		oldSetting.DeskripsiSekolah == setting.DeskripsiSekolah &&
		oldSetting.VisiSekolah == setting.VisiSekolah &&
		oldSetting.FooterSekolah == setting.FooterSekolah &&
		oldSetting.ProfilSekolah == setting.ProfilSekolah &&
		oldSetting.Facebook == setting.Facebook &&
		oldSetting.Twitter == setting.Twitter &&
		oldSetting.Instagram == setting.Instagram &&
		oldSetting.Youtube == setting.Youtube &&
		oldSetting.PetaSekolah == setting.PetaSekolah &&
		oldSetting.LogoSekolah == setting.LogoSekolah &&
		oldSetting.HeroArea1 == setting.HeroArea1 &&
		oldSetting.HeroArea2 == setting.HeroArea2 &&
		oldSetting.HeroArea3 == setting.HeroArea3 &&
		oldSetting.HeroArea4 == setting.HeroArea4 {

		return false, nil
	}

	// =====================================
	// Update
	// =====================================

	_, err = database.DB.Exec(
		`UPDATE setting SET
		id_pengguna = ?,
		nama_sekolah = ?,
		alamat_sekolah = ?,
		npsn = ?,
		nsss = ?,
		nsp = ?,
		kelurahan_sekolah = ?,
		kecamatan_sekolah = ?,
		kabupaten_sekolah = ?,
		provinsi_sekolah = ?,
		kodepos_sekolah = ?,
		telp_sekolah = ?,
		email_sekolah = ?,
		web_sekolah = ?,
		deskripsi_sekolah = ?,
		visi_sekolah = ?,
		footer_sekolah = ?,
		profil_sekolah = ?,
		facebook = ?,
		twitter = ?,
		instagram = ?,
		youtube = ?,
		peta_sekolah = ?,
		logo_sekolah = ?,
		heroarea1 = ?,
		heroarea2 = ?,
		heroarea3 = ?,
		heroarea4 = ?
	WHERE id_sekolah = ?`,
		setting.IDPengguna,
		setting.NamaSekolah,
		setting.AlamatSekolah,
		setting.NPSN,
		setting.NSSS,
		setting.NSP,
		setting.KelurahanSekolah,
		setting.KecamatanSekolah,
		setting.KabupatenSekolah,
		setting.ProvinsiSekolah,
		setting.KodePosSekolah,
		setting.TelpSekolah,
		setting.EmailSekolah,
		setting.WebSekolah,
		setting.DeskripsiSekolah,
		setting.VisiSekolah,
		setting.FooterSekolah,
		setting.ProfilSekolah,
		setting.Facebook,
		setting.Twitter,
		setting.Instagram,
		setting.Youtube,
		setting.PetaSekolah,
		setting.LogoSekolah,
		setting.HeroArea1,
		setting.HeroArea2,
		setting.HeroArea3,
		setting.HeroArea4,
		setting.IDSekolah,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}
