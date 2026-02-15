# Hướng dẫn Đóng góp (Contributing Guide)

Cảm ơn bạn quan tâm đến việc đóng góp cho dự án AI Analytics! 🎉

## 🤝 Cách đóng góp

### Báo cáo Bug

1. Kiểm tra xem bug đã được report chưa trong [Issues](https://github.com/your-repo/issues)
2. Tạo issue mới với template:

```markdown
**Mô tả bug**
Mô tả rõ ràng về bug

**Các bước để reproduce**
1. Làm A
2. Làm B
3. Thấy lỗi

**Kết quả mong đợi**
Mô tả kết quả bạn mong đợi

**Screenshots**
Nếu có, thêm screenshots

**Environment:**
- OS: [e.g. Windows 11]
- Docker version:
- Browser: [e.g. Chrome 120]
```

### Đề xuất Feature mới

1. Mở issue với label `enhancement`
2. Mô tả rõ:
   - Vấn đề bạn đang gặp
   - Giải pháp đề xuất
   - Các alternatives bạn đã cân nhắc

### Pull Request Process

1. **Fork repository**
   ```bash
   git clone https://github.com/your-username/AI_analysis.git
   cd AI_analysis
   ```

2. **Tạo branch mới**
   ```bash
   git checkout -b feature/your-feature-name
   # hoặc
   git checkout -b fix/bug-description
   ```

3. **Development**
   ```bash
   # Setup development environment
   ./manage.sh setup
   
   # Make your changes
   # ...
   
   # Test your changes
   ./manage.sh health
   ```

4. **Commit changes**
   ```bash
   git add .
   git commit -m "feat: add revenue prediction confidence interval"
   
   # Commit message format:
   # feat: new feature
   # fix: bug fix
   # docs: documentation changes
   # style: formatting, missing semi colons, etc
   # refactor: code refactoring
   # test: adding tests
   # chore: maintenance
   ```

5. **Push và tạo Pull Request**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **PR Requirements**
   - [ ] Code follows project style
   - [ ] Tests pass
   - [ ] Documentation updated
   - [ ] PR description is clear

## 📝 Coding Standards

### Go Code Style

```go
// Good
func CalculateRevenue(ctx context.Context, restaurantID string) (float64, error) {
    // Implementation
}

// Bad
func calculate_revenue(ctx context.Context, restaurant_id string) (float64, error) {
    // Implementation
}
```

- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Add comments for exported functions

### Python Code Style

```python
# Good
def calculate_features(df: pd.DataFrame) -> pd.DataFrame:
    """Calculate monthly revenue features.
    
    Args:
        df: DataFrame with raw data
        
    Returns:
        DataFrame with calculated features
    """
    pass

# Bad
def calcFeatures(df):
    pass
```

- Follow PEP 8
- Use type hints
- Add docstrings for functions
- Run `black` for formatting

### JavaScript/React Code Style

```javascript
// Good
const RevenueChart = ({ data }) => {
  const [loading, setLoading] = useState(false);
  
  return (
    <div className="chart-container">
      <ReactECharts option={option} />
    </div>
  );
};

// Bad
const revenueChart = (props) => {
  var loading = false
  return <div><ReactECharts option={props.option} /></div>
}
```

- Use functional components with hooks
- Use const/let, not var
- Add PropTypes or TypeScript types

## 🧪 Testing

### Backend Tests (Go)

```bash
cd backend
go test ./...

# With coverage
go test -cover ./...
```

### ML Tests (Python)

```bash
cd ml
python -m pytest tests/

# With coverage
python -m pytest --cov=. tests/
```

### Frontend Tests (React)

```bash
cd client
npm test

# With coverage
npm test -- --coverage
```

## 📂 Project Structure

```
AI_analysis/
├── backend/          # Go API Server
│   ├── cmd/          # Main applications
│   ├── internal/     # Private application code
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── models/
│   │   └── services/
│   └── pkg/          # Public libraries
├── etl/              # Go ETL Workers
│   ├── cmd/
│   └── internal/
├── ml/               # Python ML Services
│   ├── training/
│   ├── prediction/
│   └── tests/
├── client/           # React Frontend
│   └── src/
│       ├── api/
│       ├── components/
│       └── App.jsx
└── docs/             # Documentation
```

## 🎯 Development Workflow

### 1. Setup Development Environment

```bash
./manage.sh setup
```

### 2. Start Services

```bash
# Start all
./manage.sh start

# Or start individual services
docker-compose up -d mongodb redis
cd backend && go run cmd/api/main.go
cd client && npm run dev
```

### 3. Make Changes

- Backend: Edit Go files in `backend/`
- Frontend: Edit React files in `client/src/`
- ML: Edit Python files in `ml/`

### 4. Test Changes

```bash
# Health check
./manage.sh health

# Test API
curl http://localhost:8080/api/v1/analytics/dashboard?restaurant_id=REST001

# Check frontend
open http://localhost:3000
```

### 5. Run Tests

```bash
# Backend
cd backend && go test ./...

# Frontend
cd client && npm test

# ML
cd ml && python -m pytest
```

## 🐛 Debugging

### Backend API Debugging

```bash
# View logs
docker-compose logs -f backend

# Attach debugger (Delve)
cd backend
dlv debug cmd/api/main.go
```

### Frontend Debugging

```bash
# Chrome DevTools
# React DevTools extension

# View console errors
npm run dev
# Then open http://localhost:3000 and check Console
```

### ML Debugging

```bash
# IPython for interactive debugging
cd ml
ipython

# Jupyter for notebooks
jupyter notebook notebooks/
```

## 📚 Documentation

### Updating Documentation

When adding features, update:

1. **README.md** - High-level overview
2. **docs/api.md** - API endpoints
3. **docs/usage.md** - Usage instructions
4. **Code comments** - Inline documentation

### Documentation Style

```markdown
# Use clear headings

## Prerequisites
- Item 1
- Item 2

## Example
\```bash
command example
\```

## Result
Expected output
```

## 🎨 UI/UX Guidelines

### Design Principles

- **Simplicity**: Keep UI clean and intuitive
- **Consistency**: Use consistent colors, fonts, spacing
- **Responsiveness**: Works on mobile, tablet, desktop
- **Accessibility**: WCAG 2.1 AA compliance

### Color Palette

```css
Primary: #667eea (Purple)
Success: #10b981 (Green)
Warning: #f59e0b (Orange)
Error: #ef4444 (Red)
Background: #f5f7fa
Text: #333333
```

## 🔐 Security

### Security Best Practices

- Never commit credentials or secrets
- Use environment variables for config
- Validate all user inputs
- Use parameterized queries
- Keep dependencies updated

### Reporting Security Issues

Email security issues to: security@example.com

DO NOT create public issues for security vulnerabilities.

## 📞 Getting Help

- **Documentation**: Check [docs/](docs/)
- **Issues**: Search existing issues
- **Discussions**: Start a discussion
- **Email**: contact@example.com

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

## 🙏 Thank You!

Your contributions make this project better. Thank you! 🚀
