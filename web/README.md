# Web UI for Order Packs Calculator

Simple web interface for order packs calculator.

## Features

- ✅ Dynamic pack size management (add/remove)
- ✅ Client-side validation (positive integers)
- ✅ Sending requests to `/packs/solve` via fetch API
- ✅ Displaying results in table
- ✅ Summary information (total packs, overage)
- ✅ Friendly error messages
- ✅ Auto-fill edge-case example (263 items)

## Usage

### Starting the server

```bash
# From project root
make run

# Or directly
go run cmd/api/main.go
```

### Opening in Browser

Open http://localhost:8080 in browser.

## Interface

### 1. Pack Sizes

- Initial 3 fields for entering pack sizes
- Button **"+ Add Pack Size"** - add new field
- Button **"Remove"** - remove field (minimum one field must remain)
- Button **"Submit pack sizes change"** - save pack sizes
- Link **"Load example"** - load edge-case example

### 2. Calculate packs for order

- Field **"Items"** - number of items to pack
- Button **"Calculate"** - perform calculation

### 3. Results

Results table:
- **Pack** - pack size
- **Quantity** - number of packs

Summary information:
- **Total Packs** - total number of packs
- **Items Ordered** - items ordered
- **Items Delivered** - items delivered
- **Overage** - overage

## Validation

### Pack Sizes
- Must be positive integers
- Range: 1 - 1,000,000
- Must not duplicate
- Minimum one value

### Items (Amount)
- Must be positive integer
- Range: 1 - 1,000,000,000

## Usage Examples

### Basic Example

1. Enter pack sizes: `250`, `500`, `1000`
2. Click **"Submit pack sizes change"**
3. Enter number of items: `251`
4. Click **"Calculate"**
5. Result: 1 pack of 500 (overage: 249)

### Edge-case Example (263 items)

1. Click on link **"Load example (edge case: 263 items)"**
2. Sizes will be automatically loaded: `[250, 500, 1000, 2000, 5000]`
3. Number of items: `263`
4. Click **"Calculate"**
5. Result: optimal pack combination

## Error Handling

### Client-side Errors
- Red border around invalid field
- Error message at top of page
- Automatic disappearance after 5 seconds

### Server-side Errors
- Displaying message from server
- Error details in browser console

## Technical Details

### API Request

```javascript
fetch('/packs/solve', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({
        sizes: [250, 500, 1000],
        amount: 251
    })
})
```

### API Response

```json
{
  "solution": {
    "500": 1
  },
  "overage": 249,
  "packs": 1
}
```

## Styling

- Modern minimalist design
- Responsive layout
- Smooth transitions and animations
- Color scheme:
  - Success: green (#4CAF50)
  - Error: red (#f44336)
  - Info: blue (#2196F3)
  - Delete: red (#f44336)

## Browser Compatibility

- Chrome/Edge (latest versions)
- Firefox (latest versions)
- Safari (latest versions)

Requires support:
- ES6+ (async/await, fetch API)
- CSS Grid/Flexbox

