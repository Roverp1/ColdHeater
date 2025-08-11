# TODO's

## Daniel

1. Create personal git branch
2. Implement sender module

   1. See `internal/sender/interfaces.go` for required functions
   2. Start with `SendColdEmail()` - use Resend api
      - reference `sendEmails.ts` or `src/services/warmup/stages/resendMessage.ts` from prototype as needed
      - define structs as needed
        - reference `app_architecture.md` for databese structs like Sender (database structs should be defined in `internal/database/models.go`)
        - in case mf wants to work with database: mf can implement tables as he needs them (e.g. sender, conversation...) - seek help for work with database

3. Success criteria for SendColdEmail():
   - Sends email via Resend API
   - Returns error if sending fails
   - Basic test that verifies email sending work (preferable but optional)
