# This is just playground

i just need to know how image upload work in go

To upload an image using **Insomnia** or **Postman** to the RESTful API you built (with the `/images` endpoint), you need to send a `POST` request with a `multipart/form-data` payload, including the image file. Below are step-by-step instructions for both tools to test the image upload functionality of your Go API (running on `http://localhost:8080`).

### Prerequisites

- Your Go API (`image-upload-api`) is running locally (`go run cmd/main.go`).
- The server is accessible at `http://localhost:8080`.
- You have an image file (e.g., `image.jpg`) to upload.
- The `uploads/` directory exists in your project root (it’s created automatically by the `LocalStorage` initialization).

### Using Postman

Postman is a popular API testing tool that makes it easy to send multipart form-data requests.

1. **Open Postman**:

   - Launch Postman and create a new request.

2. **Configure the Request**:

   - **Method**: Select `POST`.
   - **URL**: Enter `http://localhost:8080/images`.

3. **Set Up the Body**:

   - Go to the **Body** tab.
   - Select `form-data` as the body type.
   - Add a new key-value pair:
     - **Key**: `image` (this matches the `FormFile("image")` in your `image_handler.go`).
     - **Type**: Change the type to `File` (click the dropdown next to the key).
     - **Value**: Click `Select Files` and choose your image file (e.g., `image.jpg`).

4. **Send the Request**:

   - Click the **Send** button.
   - If successful, you should see a `201 Created` status and a response body like:
     ```
     Image uploaded successfully: <uuid>-image.jpg
     ```

5. **Verify**:

   - Check the `uploads/` directory in your project root for the saved image (e.g., `<uuid>-image.jpg`).
   - Open `images.db` with a SQLite client and query:
     ```sql
     SELECT * FROM images;
     ```
     You should see a record with the image’s ID, filename, and path.

6. **Troubleshooting**:
   - **400 Bad Request**: Ensure the key is `image` and a file is selected.
   - **500 Internal Server Error**: Check the server logs for issues (e.g., file write permissions, database errors).

this is using bash:

```bash
curl -X POST -F "image=@/path/to/image.jpg" http://localhost:8080/images
```
