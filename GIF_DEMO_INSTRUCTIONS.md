# Creating Web UI GIF Demo

## Tools

**Recommended:**
- **macOS**: [Kap](https://getkap.co/)
- **Windows**: [ScreenToGif](https://www.screentogif.com/)
- **Linux**: [Peek](https://github.com/phw/peek)
- **Cross-platform**: [LICEcap](https://www.cockos.com/licecap/)

## Quick Start

### 1. Start application
```bash
make run
# or
make compose-up
```

### 2. Open browser
```
http://localhost:8080
```

### 3. Record scenario (~12 seconds)

1. Show initial state (2 sec)
2. Click "Load example" (1 sec)
3. Click "Calculate" (1 sec)
4. Show results (3 sec)
5. Change amount to 12001 (2 sec)
6. Show new results (2 sec)

### 4. Recording settings

- **FPS**: 15-20
- **Resolution**: 1280x720
- **Size**: < 5MB

### 5. Optimization

```bash
# Install gifsicle
brew install gifsicle  # macOS

# Optimize
gifsicle -O3 --lossy=80 --colors 128 input.gif -o output.gif
```

## Adding to README

```bash
mkdir -p docs/images
cp demo.gif docs/images/web-ui-demo.gif
git add docs/images/web-ui-demo.gif
```

In `README.md` replace:
```markdown
<!-- TODO: Add GIF demo here -->
```

With:
```markdown
![Web UI Demo](docs/images/web-ui-demo.gif)
```

## Tips

- Clean browser window (no toolbars)
- Zoom 100%
- Smooth cursor movements
- Size < 5MB for fast loading

