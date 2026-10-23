-- phpMyAdmin SQL Dump
-- version 5.2.3
-- https://www.phpmyadmin.net/
--
-- Host: localhost:3306
-- Generation Time: Sep 27, 2026 at 03:11 PM
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
-- Table structure for table `tbl_mapel`
--

CREATE TABLE `tbl_mapel` (
  `id_mapel` int NOT NULL,
  `id_induk` int DEFAULT NULL,
  `nama_mapel` varchar(100) COLLATE utf8mb4_general_ci NOT NULL,
  `urutan` int NOT NULL DEFAULT '0',
  `aktif` tinyint NOT NULL DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `tbl_mapel`
--

INSERT INTO `tbl_mapel` (`id_mapel`, `id_induk`, `nama_mapel`, `urutan`, `aktif`) VALUES
(1, NULL, 'Pendidikan Agama dan Budi Pekerti', 1, 1),
(2, NULL, 'Pendidikan Pancasila', 2, 1),
(3, NULL, 'Bahasa Indonesia', 3, 1),
(4, NULL, 'Matematika', 4, 1),
(5, NULL, 'Ilmu Pengetahuan Alam', 5, 1),
(6, NULL, 'Ilmu Pengetahuan Sosial', 6, 1),
(7, NULL, 'Bahasa Inggris', 7, 1),
(8, NULL, 'PJOK', 8, 1),
(9, NULL, 'Teknologi Informasi & Komunikasi', 9, 1),
(10, NULL, 'Seni Budaya & Prakarya', 10, 1),
(11, 10, 'Seni Musik', 1, 1),
(12, 10, 'Seni Tari', 2, 1),
(13, 10, 'Seni Rupa', 3, 1),
(14, 10, 'Seni Teater', 4, 1),
(15, 10, 'Prakarya', 5, 1),
(16, NULL, 'Proyek Penguatan Profil Pelajar Pancasila (P5)', 11, 1),
(17, NULL, 'Bimbingan Konseling', 12, 1),
(18, 16, 'Proyek 1', 1, 1),
(19, 16, 'Proyek 2', 2, 1),
(20, 16, 'Proyek 3', 3, 1),
(21, 16, 'Proyek 4', 4, 1),
(22, NULL, 'Muatan Lokal', 13, 1),
(23, 22, 'Bahasa Daerah', 1, 1);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `tbl_mapel`
--
ALTER TABLE `tbl_mapel`
  ADD PRIMARY KEY (`id_mapel`),
  ADD KEY `idx_induk` (`id_induk`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `tbl_mapel`
--
ALTER TABLE `tbl_mapel`
  MODIFY `id_mapel` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=24;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `tbl_mapel`
--
ALTER TABLE `tbl_mapel`
  ADD CONSTRAINT `fk_mapel_induk` FOREIGN KEY (`id_induk`) REFERENCES `tbl_mapel` (`id_mapel`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
