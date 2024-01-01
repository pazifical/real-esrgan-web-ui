document.getElementById("upload-file").onchange = async (event) => {
  console.log("upload file");
  const form = document.getElementById("upload-form");
  const formData = new FormData(document.getElementById("upload-form"));
  formData.append(
    "factor",
    Number(document.getElementById("upscale-factor").value)
  );
  const response = await fetch("/api/upscale", {
    method: "POST",
    body: formData,
  });
  if (response.ok) {
    event.target.value = "";
  }

  event.preventDefault();
};

document.getElementById("delete-all").onclick = async () => {
  const response = await fetch("/api/upscale", {
    method: "DELETE",
  });
  console.log(response);
};

async function updateList() {
  const response = await fetch("/api/upscale");
  const images = await response.json();

  document.getElementById("gallery").innerHTML = images
    .map(
      (i) => `<div>
        <a href="${i.filepath}" target="_blank">
            <img alt="${i.name} " src="${i.filepath}" height="150" class="gallery__img"/>
        </a> 
        </div>`
    )
    .join("");
}

setInterval(updateList, 2000);
