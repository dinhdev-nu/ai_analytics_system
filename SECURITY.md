# Security Policy

## Supported Versions

Security updates are provided for the following versions:

| Version | Supported          | End of Life |
| ------- | ------------------ | ----------- |
| 1.0.x   | :white_check_mark: | TBD         |
| < 1.0   | :x:                | 2026-02-15  |

## Reporting a Vulnerability

We take the security of AI Analytics seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Where to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to:

**security@aianalytics.example.com**

### What to Include

Please include the following information in your report:

1. **Description** - Detailed description of the vulnerability
2. **Impact** - What an attacker could achieve
3. **Steps to Reproduce** - Clear steps to reproduce the issue
4. **Proof of Concept** - PoC code or screenshots (if applicable)
5. **Affected Components** - Which parts of the system are affected
6. **Suggested Fix** - If you have ideas on how to fix it (optional)
7. **Your Contact Information** - How we can reach you

### Example Report

```
Subject: [SECURITY] SQL Injection in Analytics Endpoint

Description:
The /api/v1/analytics/dashboard endpoint is vulnerable to NoSQL injection
through the restaurant_id parameter.

Impact:
An attacker can access data from other restaurants by manipulating the
restaurant_id parameter with MongoDB operators.

Steps to Reproduce:
1. Send request: GET /api/v1/analytics/dashboard?restaurant_id[$ne]=REST001
2. Observe that data from all restaurants is returned

Proof of Concept:
curl 'http://localhost:8080/api/v1/analytics/dashboard?restaurant_id[$ne]=REST001'

Affected Components:
- backend/internal/handlers/analytics_handler.go (line 45)
- backend/internal/services/analytics_service.go (line 120)

Suggested Fix:
Validate and sanitize the restaurant_id parameter to only accept
alphanumeric characters.

Contact: john.doe@example.com
```

## Response Timeline

We aim to respond to security reports according to the following timeline:

- **Initial Response:** Within 24 hours
- **Validation:** Within 3 business days
- **Status Update:** Every 5 business days
- **Fix Released:** Depends on severity (see below)

### Severity Levels

| Severity | Response Time | Fix Release |
|----------|--------------|-------------|
| **Critical** | < 24 hours | 1-3 days |
| **High** | < 48 hours | 1-2 weeks |
| **Medium** | < 1 week | 2-4 weeks |
| **Low** | < 2 weeks | Next release |

### Severity Criteria

**Critical:**
- Remote code execution
- Authentication bypass
- Data breach affecting all users
- Privilege escalation to admin

**High:**
- SQL/NoSQL injection
- XSS allowing account takeover
- Exposure of sensitive data
- Authentication flaws

**Medium:**
- CSRF vulnerabilities
- Information disclosure (limited scope)
- Denial of Service
- Weak cryptography

**Low:**
- Non-sensitive information disclosure
- Minor security misconfigurations
- Security best practice violations

## Disclosure Policy

### Coordinated Disclosure

We follow a coordinated disclosure process:

1. **You report** the vulnerability privately
2. **We acknowledge** receipt within 24 hours
3. **We investigate** and develop a fix
4. **We release** the fix in a security update
5. **We publish** a security advisory 7-14 days after the fix
6. **You can disclose** publicly after the advisory (if desired)

### Public Disclosure

We believe in transparency and will publicly disclose:

- All security vulnerabilities after they are fixed
- Credit to the reporter (if desired)
- Technical details and impact assessment
- Mitigation steps for users

Public disclosures will be made through:
- GitHub Security Advisories
- Release notes
- Security page (SECURITY.md)
- Project blog/newsletter

## Security Updates

### How to Stay Informed

Subscribe to security updates:

1. **Watch** this repository on GitHub (Releases only)
2. **Subscribe** to our security mailing list: security-announce@aianalytics.example.com
3. **Follow** our [Security Advisories](../../security/advisories)

### Applying Security Updates

When a security update is released:

1. **Backup your data**:
   ```bash
   ./manage.sh backup
   ```

2. **Pull latest changes**:
   ```bash
   git pull origin main
   ```

3. **Rebuild containers**:
   ```bash
   docker-compose down
   docker-compose build
   docker-compose up -d
   ```

4. **Verify the fix**:
   ```bash
   ./manage.sh health
   ```

## Security Best Practices

### For Users

#### Production Deployment

1. **Enable Authentication** (when available in v1.1.0)
2. **Use Strong Passwords**
   - MongoDB: 16+ characters, mixed case, numbers, symbols
   - Redis: AUTH enabled with strong password
3. **Enable HTTPS**
   - Use Let's Encrypt for free SSL certificates
   - Enforce HTTPS redirects
4. **Configure Firewall**
   - Only expose necessary ports (443, 80)
   - Restrict MongoDB/Redis to localhost
5. **Regular Backups**
   ```bash
   # Daily automated backups
   0 3 * * * /path/to/manage.sh backup
   ```
6. **Keep Updated**
   - Apply security patches promptly
   - Subscribe to security announcements
7. **Monitor Logs**
   - Check for suspicious activity
   - Set up alerts for errors
8. **Use Secrets Management**
   - Don't commit secrets to Git
   - Use environment variables or vault

#### Environment Variables

Never commit these to Git:
```env
# .env (add to .gitignore)
MONGODB_URI=mongodb://user:STRONG_PASSWORD@localhost:27017
REDIS_PASSWORD=ANOTHER_STRONG_PASSWORD
JWT_SECRET=RANDOM_64_CHARACTER_STRING
API_KEY=SECRET_API_KEY
```

#### Network Security

```yaml
# docker-compose.yml
services:
  mongodb:
    networks:
      - backend
    # Don't expose ports in production
    # ports:
    #   - "27017:27017"
  
  backend:
    networks:
      - backend
      - frontend
    ports:
      - "8080:8080"
      
networks:
  backend:
    internal: true
  frontend:
```

### For Contributors

#### Secure Coding Guidelines

1. **Input Validation**
   ```go
   // Bad
   restaurantID := c.Query("restaurant_id")
   
   // Good
   restaurantID := c.Query("restaurant_id")
   if !regexp.MustCompile(`^[A-Z0-9]+$`).MatchString(restaurantID) {
       c.JSON(400, gin.H{"error": "Invalid restaurant_id"})
       return
   }
   ```

2. **SQL/NoSQL Injection Prevention**
   ```go
   // Bad
   filter := bson.M{"restaurant_id": restaurantID}
   
   // Good - parametrized query
   filter := bson.M{
       "restaurant_id": bson.M{"$eq": restaurantID},
   }
   ```

3. **XSS Prevention**
   ```javascript
   // Bad
   element.innerHTML = userInput;
   
   // Good
   element.textContent = userInput;
   // Or use React (auto-escapes by default)
   ```

4. **Authentication** (when implemented)
   ```go
   // Use bcrypt for password hashing
   hashedPassword, err := bcrypt.GenerateFromPassword(
       []byte(password), 
       bcrypt.DefaultCost,
   )
   ```

5. **Rate Limiting**
   ```go
   // Already implemented in backend
   router.Use(middleware.RateLimiter(100, time.Minute))
   ```

#### Security Review Checklist

Before submitting a PR:

- [ ] No secrets in code or commits
- [ ] Input validation on all user inputs
- [ ] Parameterized database queries
- [ ] Error messages don't leak sensitive info
- [ ] Authentication/authorization checked (when applicable)
- [ ] HTTPS used for external APIs
- [ ] Dependencies are up to date
- [ ] No known vulnerabilities in dependencies

#### Dependency Management

```bash
# Go
go get -u ./...
go mod tidy

# Python
pip install --upgrade -r requirements.txt
pip-audit  # Check for vulnerabilities

# Node.js
npm audit
npm audit fix
```

## Known Vulnerabilities

### Current (v1.0.0)

**CVE-NONE:** No known vulnerabilities

### Acknowledged Limitations

1. **No Authentication** (v1.0.0)
   - All API endpoints are public
   - Planned for v1.1.0
   - **Mitigation:** Use firewall rules to restrict access

2. **No Rate Limiting per User** (v1.0.0)
   - Rate limiting is per-IP only
   - Planned for v1.1.0
   - **Mitigation:** Configure strict IP-based rate limits

3. **MongoDB Without Auth** (Default Development Setup)
   - Development docker-compose has no auth
   - **Mitigation:** Enable in production (see docs/deployment.md)

## Security Tools

We use the following tools to maintain security:

- **Dependabot:** Automated dependency updates
- **CodeQL:** Static analysis for vulnerabilities
- **gosec:** Go security checker
- **bandit:** Python security checker
- **npm audit:** Node.js dependency vulnerabilities
- **Docker Scan:** Container vulnerability scanning

## Bug Bounty Program

**Status:** Not available yet

We are considering a bug bounty program for future versions. Stay tuned!

## Security Hall of Fame

We recognize security researchers who help us:

| Researcher | Vulnerability | Severity | Date |
|------------|--------------|----------|------|
| *TBD*      | *TBD*        | *TBD*    | *TBD* |

## Contact

### Security Team

- **Email:** security@aianalytics.example.com
- **PGP Key:** [Download PGP Key](#)
- **Response Time:** < 24 hours

### General Inquiries

- **GitHub Issues:** [Report Non-Security Bugs](#)
- **Email:** contact@aianalytics.example.com

## Legal

### Safe Harbor

We support safe harbor for security researchers who:

- Make a good faith effort to avoid privacy violations and data destruction
- Only interact with accounts you own or with explicit permission
- Do not exploit vulnerabilities beyond Proof of Concept
- Report vulnerabilities promptly
- Keep vulnerability details confidential until we release a fix

We will not pursue legal action against researchers who follow these guidelines.

### Acknowledgment

This security policy is inspired by:
- [GitHub's Security Policy](https://github.com/github/docs/security/policy)
- [Node.js Security Policy](https://github.com/nodejs/node/security/policy)
- [OWASP Guidelines](https://owasp.org/)

---

**Version:** 1.0  
**Last Updated:** February 15, 2026

**Thank you for keeping AI Analytics secure!**
