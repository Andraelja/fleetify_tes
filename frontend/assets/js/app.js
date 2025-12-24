const API_BASE_URL = "http://localhost:8080";

function getToken() {
  return localStorage.getItem("token");
}

function apiRequest({ url, method = "GET", data = null, success, error }) {
  $.ajax({
    url: API_BASE_URL + url,
    method: method,
    contentType: "application/json",
    headers: {
      Authorization: "Bearer " + getToken(),
    },
    data: data ? JSON.stringify(data) : null,
    success: success,
    error: function (xhr) {
      let message = xhr.responseJSON?.message || "Terjadi kesalahan server";

      Swal.fire({
        icon: "error",
        title: "Error",
        text: message,
      });

      if (error) error(xhr);
    },
  });
}
