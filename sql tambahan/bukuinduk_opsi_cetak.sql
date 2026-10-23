-- phpMyAdmin SQL Dump
-- version 5.2.3
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Sep 27, 2026 at 03:08 PM
-- Server version: 8.4.3
-- PHP Version: 8.3.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `smpn254`
--

-- --------------------------------------------------------

--
-- Table structure for table `bukuinduk_opsi_cetak`
--

CREATE TABLE `bukuinduk_opsi_cetak` (
  `id` tinyint NOT NULL DEFAULT '1',
  `sampul` tinyint NOT NULL DEFAULT '1',
  `daftar_isi` tinyint NOT NULL DEFAULT '1',
  `pemeriksaan` tinyint NOT NULL DEFAULT '1',
  `nomor_urut` tinyint NOT NULL DEFAULT '1',
  `petunjuk` tinyint NOT NULL DEFAULT '1',
  `lembar_iv` tinyint NOT NULL DEFAULT '1',
  `halaman_akhir` tinyint NOT NULL DEFAULT '1',
  `tingkat` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '7,8,9',
  `kepsek_nama_override` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `kepsek_nip_override` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `walikelas_nama_override` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `walikelas_nip_override` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `bukuinduk_opsi_cetak`
--

INSERT INTO `bukuinduk_opsi_cetak` (`id`, `sampul`, `daftar_isi`, `pemeriksaan`, `nomor_urut`, `petunjuk`, `lembar_iv`, `halaman_akhir`, `tingkat`, `kepsek_nama_override`, `kepsek_nip_override`, `walikelas_nama_override`, `walikelas_nip_override`) VALUES
(1, 1, 1, 1, 1, 1, 1, 0, '', '', '', '', '');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `bukuinduk_opsi_cetak`
--
ALTER TABLE `bukuinduk_opsi_cetak`
  ADD PRIMARY KEY (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
