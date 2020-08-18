function jumpToResult() {
	// credit: https://stackoverflow.com/questions/13735912/anchor-jumping-by-using-javascript
	var url = location.href;
	location.href = '#result-header';
	history.replaceState(null, null, url);
}
