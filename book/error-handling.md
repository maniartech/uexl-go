# Error Handling

UExL provides clear and actionable error information when expression evaluation fails.

## Error Object
When an error occurs, UExL returns an object with:
- **errorCode**: Standardized code (e.g., `SYNTAX_ERROR`, `TYPE_MISMATCH`, `RUNTIME_ERROR`)
- **message**: Descriptive message
- **line** and **column**: Location in the expression

## Example
```
{
  "errorCode": "SYNTAX_ERROR",
  "message": "Unexpected token '+' encountered.",
  "line": 3,
  "column": 15
}
```

## Debugging Support
- **Debug Mode**: Logs detailed evaluation steps and intermediate values.
- **Stack Traces**: Provided for runtime errors in nested expressions or function calls.