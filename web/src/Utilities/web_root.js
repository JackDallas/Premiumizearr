export function CalculateAPIPath(path) {
  let webRoot = window.location.href;

  if (webRoot.indexOf("index.html") > -1) {
    webRoot = webRoot.substring(0, webRoot.indexOf("index.html"));
  }

  if (webRoot[webRoot.length - 1] !== "/") {
    webRoot += "/";
  }

  if (path[0] == "/") {
    // console.log(webRoot + path.substring(1));
    return webRoot + path.substring(1);
  }
  
  // console.log(webRoot + path);

  return webRoot + path;
}

