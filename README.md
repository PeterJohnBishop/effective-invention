# effective-invention

## JS 

Sends batches of 10 requests at a time with a brief wait to stay under a 1000 request per minut rate limit.

Result: fetched 69726 tasks in 268.69 seconds. RMP: 156.31

## Go

Sends a new request whenever a token is availible, with the rate limited by 1 token added every 60 milliseconds to stay under the 1000 request per minute rate limit. Because each request is in its own go routine it doesn't wait for other requests to finish. Responses are sent through a buffered channel into a processor that handles combining the tasks in a separate go routine. 

Result: fetched 69726 tasks in 44.543s seconds. RPM: 939.23

# High-Trust Digital Vault

- secure file storage and sharing via time limited links embedded in QR codes, sent directly via email and requiring biometric verification to view.


## Ideas

1. The "Secure VIP" Event Management System

Instead of simple tickets, use facial comparison to ensure the person holding the QR code is the actual ticket owner.

The Workflow: A user registers and uploads a "master" selfie (stored via your User Management/S3). They receive a temporary, expiring QR code via your custom email API.

The Security: At the door, the staff scans the QR code. The app then captures a live photo of the guest and uses AWS Rekognition (Facial Comparison) to match it against the master image.

Expansion: Add a "Press Kit" feature where approved journalists can download high-res assets via expiring links once their identity is verified.

2. "No-Password" Sensitive Document Sharing

A B2B platform for lawyers, real estate agents, or accountants to share sensitive PDFs without relying on passwords that can be leaked.

The Workflow: The sender uploads a document (closing papers, tax returns). They send a link via your custom email API.

The Security: To open the link, the recipient must perform a "Facial Analysis" check via their webcam to verify they are the person on file.

The "Expiring" Twist: Once the face is verified, the download link becomes active for only 60 seconds, then expires completely to prevent link-sharing.

3. Smart Asset Management for Retail/Logistics

A system to track high-value physical assets (like high-end rental equipment or secure inventory) using QR codes and human accountability.

The Workflow: Every piece of equipment has a permanent QR code. When a worker "checks out" an item, they scan the code, which triggers a file upload of the item’s current condition photo.

Facial Analysis: Use Rekognition to ensure the worker isn't wearing a mask or sunglasses during checkout (Facial Analysis for "Liveness") and verify their identity.

The Log: The system emails the manager a summary of the checkout with the "before" photo and the identity of the worker.

4. "Legacy Vault" (Time-Locked Inheritance)

A digital "In case of emergency" vault for credentials, wills, or personal messages.

The Workflow: Users upload encrypted files. They set "Check-in" intervals (e.g., once every 6 months).

The Trigger: If a user fails to check in, the system uses your email API to notify "Beneficiaries."

The Verification: To claim the files, the beneficiary must upload a photo that matches a pre-approved identity document (AWS Rekognition Comparison). Upon success, the API generates expiring download links for the vault contents.

1. The "Face-Locked" Gateway (Using your Rekognition API)

Instead of emailing the QR code as an image attachment, email a secure link to a "Verification Gateway."

The Workflow: The user clicks the link in their email. Before they see the QR code, the browser requests access to their camera.

The Security: Your API uses AWS Rekognition (Facial Comparison) to match the live caller against the user's profile photo.

The Result: The QR code only renders on the screen after the face is verified. An attacker who steals the email cannot pass the biometric check.

2. "Split-Channel" Delivery (Out-of-Band)

Never send the "Key" and the "Lock" in the same place.

The Workflow: Send the email with the QR code, but password-protect the image or the landing page it points to.

The Security: Send the password (or a PIN) via a different channel, like an SMS or a Push Notification.

The Result: Even if the email is intercepted, the QR code is a "dead" link without the second factor.

3. Dynamic QR Codes with "Burn-on-Read"

If you use a static QR code, it can be screenshotted and reused. You should use your Expiring Links logic to make the QR code dynamic.

Single Use: Set the link behind the QR code to expire the moment it is successfully scanned once.

Short TTL (Time-To-Live): Generate the QR code image on-the-fly when the email is opened, and set the URL to expire in 5–10 minutes.

The Result: By the time a hacker might find the email in a "Sent" folder or a compromised archive, the link is already dead.

4. Visual Obfuscation (The "Privacy Mask")

Email providers (like Gmail) often "pre-fetch" images to scan them, which can accidentally trigger "one-time use" links.

Interactive Loading: Don't embed the QR code directly. Embed a "Click to Reveal" button.

The Security: This prevents email bots from "consuming" the single-use link before the human actually sees it. It also ensures the QR code isn't sitting visible in a notification preview on a locked phone screen.