# resize
- Image resizer backend service

## POST
- `http://localhost:8080/resize`
- Send the value for respective values and image using forms from frontend.

## Body
### form-data

| Key     | Value                                               |
|---------|-----------------------------------------------------|
| width   | 400                                                 |
| height  | 500                                                 |
| quality | 95                                                  |
| format  | png                                                 |
| file    | /C:/Users/adity/OneDrive/Pictures/Aditya_Arc.png    |