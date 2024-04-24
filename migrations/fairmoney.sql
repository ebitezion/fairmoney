-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Jan 18, 2024 at 09:11 AM
-- Server version: 10.4.27-MariaDB
-- PHP Version: 8.2.0

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `fairmoney_db`
--

-- --------------------------------------------------------

--
-- Table structure for table `permissions`
--

CREATE TABLE `permissions` (
  `id` bigint(20) NOT NULL,
  `code` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `permissions`
--

INSERT INTO `permissions` (`id`, `code`) VALUES
(1, 'account:read'),
(2, 'account:write');

-- --------------------------------------------------------

--
-- Table structure for table `tokens`
--

CREATE TABLE `tokens` (
  `hash` varbinary(255) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `expiry` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `scope` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `tokens`
--

INSERT INTO `tokens` (`hash`, `user_id`, `expiry`, `scope`) VALUES
(0x3a59ff4347eda51eccaf025fb58d71212f44119434d4a5511ae60425c261eb37, 21, '2024-01-15 11:17:53', 'activation'),
(0x3ed185881aef711a30a5b48211d898306c344e94eedf4b11873ee577e1dc80c1, 20, '2024-01-14 05:43:10', 'authentication'),
(0x4bafe08bf64f9553a2b4476a8fea3509a473627fb970c44c2ef40f088fb39bd9, 20, '2024-01-14 05:52:09', 'authentication'),
(0x54a3e02a8c8325edf05d807b5e328f599b1bb0fe9102224174dfbd1c2a836c4d, 20, '2024-01-16 07:50:24', 'authentication'),
(0x558fff28d28f6aaf580366521a96bbfce6e804b2065fb03c431c96adced2405f, 24, '2024-01-17 13:15:47', 'authentication'),
(0x86de0da06ebc206315d35e996db3c72d361bf25a145c26324cf7b0eabc48bbcc, 20, '2024-01-15 11:12:42', 'activation'),
(0xbdf99e7595e7b6d70af4b98e27e90bb31f57e947d528b5866d9e1e10b950a096, 20, '2024-01-17 08:59:31', 'authentication'),
(0xc2d1ee0606c48d0f9c4970f60805c77a02394e58b504380fa2c11ecfc6ca6757, 20, '2024-01-14 05:52:26', 'authentication');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `email` varchar(200) NOT NULL,
  `activated` tinyint(1) NOT NULL DEFAULT 0,
  `password_hash` text NOT NULL,
  `created_at` date NOT NULL DEFAULT current_timestamp(),
  `version` int(11) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `name`, `username`, `email`, `activated`, `password_hash`, `created_at`, `version`) VALUES
(1, 'Alice Smith', 'alice2023', 'alice@example.com', 1, '$2a$12$4aNd7.6Q7k.eWeXZMqSoGeue//nDZFW33FbdHbFqf.idLR.fejtWy', '2024-01-11', 1),
(2, 'Alicia Smith', 'alice2024', 'alice123@example.com', 0, '$2a$12$7wIPNd4XduoUh/oO7hflCeFD0kVL2kTq90HtvGVuWEQYCQRc2Mj3O', '2024-01-11', 1),
(20, 'Ebite Ogochukwu', 'ebitezion', 'ebitezion@gmail.com', 1, '$2a$12$c4lymVcloTbYy5zzLksuH.UjnhvWdmp9JErAkBJzsgIk6XDj.CtpC', '2024-01-12', 1),
(21, 'Ebite Ogochukwu', 'ebitezion360', 'ebitezion360@gmail.com', 0, '$2a$12$1QFMcfUD1NQDItUex6q7i.7AOBzgdREDA23YYUGJWxe1c8/NGbzI2', '2024-01-12', 1),
(25, 'David Ade', 'david360', 'akanbiadenugba699@gmail.com', 0, '$2a$12$BSvYxQTUvFcwsHxVD1Imp.cRCKq2LsnHQTXt5TGzGulrH.gVpQPsm', '2024-01-16', 1);

-- --------------------------------------------------------

--
-- Table structure for table `users_permissions`
--

CREATE TABLE `users_permissions` (
  `user_id` bigint(20) NOT NULL,
  `permission_id` bigint(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `users_permissions`
--

INSERT INTO `users_permissions` (`user_id`, `permission_id`) VALUES
(1, 1),
(2, 1),
(20, 1),
(21, 1);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `permissions`
--
ALTER TABLE `permissions`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `tokens`
--
ALTER TABLE `tokens`
  ADD PRIMARY KEY (`hash`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `username` (`username`,`email`),
  ADD UNIQUE KEY `username_2` (`username`);

--
-- Indexes for table `users_permissions`
--
ALTER TABLE `users_permissions`
  ADD PRIMARY KEY (`user_id`,`permission_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `permissions`
--
ALTER TABLE `permissions`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=26;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
