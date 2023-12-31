document.getElementById("upload-file").onchange = async (event) => {
    console.log("upload file")
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

    document.getElementById("processed").innerHTML = images.map(i => `<li>${i.name} <a href="${i.filepath}">Open</a></li>`)
}

window.setTimeout(updateList, 1000);