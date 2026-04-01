# Stockyard Roster

**CRM for solo founders.** Contacts, notes, follow-up reminders, deal pipeline. Not Salesforce. The thing you actually want before Salesforce. Single binary, no external dependencies.

Part of the [Stockyard](https://stockyard.dev) suite of self-hosted developer tools.

## Quick Start

```bash
curl -sfL https://stockyard.dev/install/roster | sh
roster
```

## Usage

```bash
# Add a contact
curl -X POST http://localhost:8930/api/contacts \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Smith","email":"jane@acme.com","company":"Acme","stage":"lead"}'

# Add a note
curl -X POST http://localhost:8930/api/contacts/{id}/activities \
  -H "Content-Type: application/json" \
  -d '{"type":"call","content":"Discussed pricing, following up Thursday"}'

# Create a deal
curl -X POST http://localhost:8930/api/deals \
  -H "Content-Type: application/json" \
  -d '{"contact_id":"{id}","title":"Enterprise license","value_cents":500000,"stage":"prospect"}'

# Set a reminder
curl -X POST http://localhost:8930/api/reminders \
  -H "Content-Type: application/json" \
  -d '{"contact_id":"{id}","content":"Send proposal","due_at":"2026-04-05"}'
```

## Free vs Pro

| Feature | Free | Pro ($2.99/mo) |
|---------|------|----------------|
| Contacts | 25 | Unlimited |
| Deal pipeline | ✓ | ✓ |
| Reminders | ✓ | ✓ |
| Activity log | ✓ | ✓ |

## License

Apache 2.0 — see [LICENSE](LICENSE).
