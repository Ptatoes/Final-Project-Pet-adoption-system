# Final Project Pet adoption system
 
Team Members

1. AINULAAQELAH BINTI HUZAIFAH (5999241101)
2. ANIS ZARIFAH BINTI TALIB (5999241105)

youtube link: https://youtu.be/3xrNA0xJiVc?si=aVGRnGHr5LVHDMP3
# Plan for the Pet Adoption Management System

## Overview

This project is a **Pet Adoption Management System** built using Go. It allows users to manage pets available for adoption, update their status, and remove them from the system. The following sections explain the page structure, user interaction flow, and backend functionality.

## Features

### Header and Navigation
![image](https://github.com/user-attachments/assets/d2cbb945-103d-4d84-a635-9f7041fb39e1)

- The page starts with a header displaying the Index page as home.
- A navigation menu (navbar) provides links to various sections of the site:
  - **Home**
  - **Manage Pets**: The current page for managing pets.
  - **Manage Adopters**: A link for managing adopters.
  - **Manage Adoptions**: A link for managing adoptions.
  - **FAQ**: A link to the Frequently Asked Questions page.

### Success Message
![image](https://github.com/user-attachments/assets/3f63544c-7cdd-4ce5-a7e8-f306042348e8)

- If a `successMessage` is available (likely passed from the Go server), it is displayed in a styled div with the class `success-message`.
- The message fades out after 1 second, and the div is removed after fading.

## Manage Pets page

### Add New Pet Form
![image](https://github.com/user-attachments/assets/19ef89fd-261e-4433-a7ba-3ea51355b4e5)


- A form allows users to add a new pet:
  - **Name** of the pet (text input).
  - **Species** of the pet (text input).
- Upon submission, the form sends a POST request to `/pets/add`, which is handled by the Go backend to add the new pet to the database.

### Available Pets Table
![image](https://github.com/user-attachments/assets/3b878970-b747-4534-a315-5cb13aef6749)

- If there are pets available (passed as `.pets`), a table is rendered showing:
  - **ID**
  - **Name**
  - **Species**
  - **Status**: Whether the pet is "Available" or "Adopted."
- Each pet has two actions:
  - **Update Status**: A dropdown to change the pet's status between "Available" and "Adopted." A form sends a POST request to `/pets/update` to update the pet's status.
  - **Delete Pet**: A form to delete the pet. A confirmation prompt appears when the user tries to delete a pet. Upon confirmation, a POST request is sent to `/pets/delete` to delete the pet from the database.

### No Pets Found
![image](https://github.com/user-attachments/assets/833dbdfb-ceef-4dd2-b980-4bb986c6bf64)

- If no pets are available, a message is displayed indicating that no pets match the search query at the adoption page.

## Manage Adopters page
![image](https://github.com/user-attachments/assets/11d8484e-b707-4296-99da-c7e5650358a9)

- **View Adopters**: Displays a list of adopters (ID, name, contact) with options to edit, delete, or add new adopters.
- **Add Adopter**: Form to submit new adopter details (name, contact) and save to the database.
- **Edit Adopter**: Allows updating an adopter's information (name, contact).
- **Delete Adopter**: Deletes an adopter after confirmation.

## Manage Adoptions page

![image](https://github.com/user-attachments/assets/6f3d2d53-78b0-4af0-bf59-dead9375b368)

   - Displays adoption history, available pets, and adopters.
   - Fetches adoption details from the database, including:
     - Pet name
     - Adopter name
     - Adoption date
   - If no pets are available for adoption, a message is displayed.

### Assign Adoption
   - Allows assigning a pet to an adopter and updating the pet's status to "Adopted."
   - After a successful adoption, the user is redirected back to the adoptions page.


### Footer
- The footer includes copyright information for the **"Pet Adoption Management System"**.

---

## Flow of the Application

### Initial Request
- The client (browser) makes a request to the Go backend for the pet management page. The Go server responds with the HTML page.


## Pet Management Page

### Displaying the Pet Management Page
- The Go server passes data to the HTML template:
  - **`successMessage`**: A potential success message.
  - **`pets`**: A list of pets containing details like `id`, `name`, `species`, and `status`.

### Adding a Pet
- When a user submits the form to add a new pet, the Go backend handles the POST request at `/pets/add`. The data (name, species) is saved to the database, and the page reloads with an updated list of pets.

### Updating Pet Status
- When a user selects a new status for a pet and clicks "Update", the form sends a POST request to `/pets/update`, and the pet's status is updated in the database.

### Deleting a Pet
- When the user clicks "Delete", a confirmation dialog appears. Upon confirmation, a POST request is sent to `/pets/delete` to remove the pet from the database.

### Success Message
- After adding, updating, or deleting a pet, the server sends back a success message, which is displayed to the user for a brief moment.

### Displaying the Updated Pet List
- After any modification (adding, updating, or deleting a pet), the page reloads and the updated list of pets is displayed.

## Adopter  Management Page

- **List Adopters**: A GET request to `/adopters` fetches and displays all adopters.
- **Add Adopter**: A POST request to `/adopters/add` saves the new adopter to the database.
- **Edit Adopter**: A GET request to `/adopters/edit/:id` fetches details for editing. A POST request updates the adopter's info.
- **Delete Adopter**: A POST request to `/adopters/delete` removes an adopter from the database.

## Adoption Management Page

- Fetch the list of adopters from the adopters table in the database.
- If the query succeeds, scan each row and store the adopter’s id, name, and contact in a slice.
- If an optional success message is passed (e.g., for adding or deleting an adopter), include it in the response.
- Render the HTML page `adopters.html` with the fetched adopter data.

---

## Backend Logic (Go Functions)

## Pet  Management Page

### Add Pet (`/pets/add`)
- Receives POST data for the pet's name and species.
- Validates and adds the new pet to the database.
- Redirects back to the pet management page with a success message.

### Update Pet Status (`/pets/update`)
- Receives POST data for the pet ID and new status.
- Updates the pet's status in the database.
- Redirects back to the pet management page with a success message.

### Delete Pet (`/pets/delete`)
- Receives POST data for the pet ID to delete.
- Deletes the pet from the database.
- Redirects back to the pet management page with a success message.

## Adopter  Management Page

- **ListAdopters**: Queries the database for adopters and renders them in the `adopters.html` template.
- **AddAdopter**: Takes form data (name, contact), inserts it into the database, and redirects with a success message.
- **DeleteAdopter**: Accepts an adopter ID, deletes the adopter from the database, and redirects with a success message.
- **ShowEditForm**: Fetches an adopter's details for editing.
- **EditAdopter**: Updates an adopter's details in the database and redirects with a success message.

## Adoption Management Page

- **ListAdopters**: Query the database for adopters and display them.
- **AddAdopter**: Insert new adopter into the database.
- **DeleteAdopter**: Remove adopter from the database.
- **ShowEditForm**: Fetch adopter data for editing.
- **EditAdopter**: Update adopter details in the database.
---

## UI and User Interaction
## Pet Management Page
The page allows the admin or user to manage pets by:
- Adding new pets to the system.
- Viewing a list of pets with their current status.
- Updating the status of each pet (to either "Available" or "Adopted").
- Deleting pets from the system.

## Adopter  Management Page

### Adopters List Page (`adopters.html`):
- Displays a table of adopters with options to edit, delete, or add adopters.
- Success messages are shown after actions like adding, updating, or deleting an adopter.

### Edit Adopter Page (`edit_adopter.html`):
- Displays the adopter’s current details in a form for editing.

## Adoption Management Page

### Adopters List Page:
- Display a table of adopters with options to add, edit, and delete.
- Show success messages upon successful actions (add, edit, delete).

### Add/Edit Adopter Page:
- Form for adding or editing adopter details.
- Upon form submission, the data is sent to the backend for processing.

### Delete Adopter:
- Option to delete adopter directly from the list.
- Confirm deletion and redirect with a success message.
