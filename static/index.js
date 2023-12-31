document.getElementById("upload-file").onchange = async (event) => {
  console.log("upload file");
  const form = document.getElementById("upload-form");
  const formData = new FormData(document.getElementById("upload-form"));
  const response = await fetch("/api/upscale", {
    method: "POST",
    body: formData,
  });
  if (response.ok) {
    document.getElementById(
      "upload-info"
    ).innerText = `Successfully uploaded ${event.target.value}`;
    event.target.value = "";
  }

  event.preventDefault();
};

async function updateList() {
  const response = await fetch("/api/upscale");
  const images = await response.json();

  document.getElementById("processed").innerHTML = images
    .toSorted((a, b) => a.name > b.name)
    .map(
      (i) => `<div>
        <a href="${i.filepath}" target="_blank">
            <img alt="${i.name} " src="${i.filepath}" height="150" />
        </a> 
        </div>`
    )
    .join("");
}

setInterval(updateList, 2000);
