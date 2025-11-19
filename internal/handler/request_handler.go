package handler

import (
	"net/http"
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>User Profile Stats Search</title>

		<style>
			body {
				font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
				max-width: 600px;
				margin: 3rem auto;
				padding: 2rem;
				text-align: center;
			}

			/* 🔵 Top GitHub banner */
			.top-banner {
				background: #f5f5f5;
				padding: 12px 20px;
				border-radius: 10px;
				margin-bottom: 2rem;
				font-size: 1rem;
				font-weight: 600;
			}

			.top-banner a {
				color: #0366d6;
				text-decoration: none;
			}

			.top-banner a:hover {
				text-decoration: underline;
			}

			h1 {
				font-size: 2rem;
				margin-bottom: 1rem;
			}

			.input-group {
				margin-bottom: 1rem;
				text-align: left;
			}

			label {
				font-weight: 600;
			}

			input {
				width: 100%;
				padding: 10px;
				margin-top: 4px;
				font-size: 1rem;
				border-radius: 8px;
				border: 1px solid #999;
			}

			button {
				width: 100%;
				padding: 12px;
				font-size: 1.1rem;
				border: none;
				border-radius: 8px;
				background-color: #007bff;
				color: white;
				cursor: pointer;
				margin-top: 1rem;
			}

			button:hover {
				background-color: #005fcc;
			}

			.note {
				margin-top: 1rem;
				opacity: .7;
				font-size: .9rem;
			}
		</style>

		<script>
			function searchStats() {
				const leetcode = document.getElementById("leetcode").value.trim();
				const codechef = document.getElementById("codechef").value.trim();
				const gfg = document.getElementById("gfg").value.trim();
				const hackerrank = document.getElementById("hackerrank").value.trim();
				const codeforces = document.getElementById("codeforces").value.trim();

				const url = 
					"/stats?"
					+ "leetcode=" + encodeURIComponent(leetcode)
					+ "&codechef=" + encodeURIComponent(codechef)
					+ "&gfg=" + encodeURIComponent(gfg)
					+ "&hackerrank=" + encodeURIComponent(hackerrank)
					+ "&codeforces=" + encodeURIComponent(codeforces);

				window.open(url, "_blank");
			}
		</script>
	</head>

	<body>
		<!-- 🔵 GitHub Project Documentation -->
		<div class="top-banner">
			📘 View Project Documentation → 
			<a href="https://github.com/mearjuntripathi/coding-profile-service" target="_blank">
				GitHub Repository
			</a>
		</div>

		<h1>Search Coding Profiles</h1>

		<div class="input-group">
			<label>LeetCode Username</label>
			<input id="leetcode" placeholder="e.g. mearjuntripathi">
		</div>

		<div class="input-group">
			<label>CodeChef Username</label>
			<input id="codechef" placeholder="e.g. isthisarjun">
		</div>

		<div class="input-group">
			<label>GeeksForGeeks Username</label>
			<input id="gfg" placeholder="e.g. mearjuntripathi">
		</div>

		<div class="input-group">
			<label>HackerRank Username</label>
			<input id="hackerrank" placeholder="e.g. mearjuntripathi">
		</div>

		<div class="input-group">
			<label>Codeforces Username</label>
			<input id="codeforces" placeholder="e.g. arjun_cf">
		</div>

		<button onclick="searchStats()">Search</button>

		<div class="note">A new tab will open showing all your stats.</div>

	</body>
	</html>
	`))
}