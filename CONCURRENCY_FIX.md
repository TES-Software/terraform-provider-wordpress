# Concurrency Fix for terraform-provider-wordpress

## Problem
The provider was experiencing crashes with "fatal error: concurrent map writes" when running `terraform plan`. This was caused by the `github.com/sogko/go-wordpress` library using the thread-unsafe `github.com/parnurzeal/gorequest` library internally.

## Root Cause
The issue occurred because:
1. Terraform runs resource operations concurrently in separate goroutines
2. All resources shared the same WordPress client instance
3. The underlying `gorequest` library is not thread-safe and has concurrent map write issues

## Solution
The fix changes the approach from sharing a single client to creating a new client for each operation:

1. **Changed `UserResource.client`** from `*wcl.Client` to `*wcl.Options` 
2. **Added `newClient()` method** that creates a fresh client for each operation
3. **Updated all CRUD operations** to use `u.newClient()` instead of the shared client
4. **Modified provider configuration** to pass client options instead of a pre-created client

## Benefits
- **Thread Safety**: Each goroutine gets its own HTTP client with isolated state
- **No Performance Impact**: Client creation is lightweight 
- **Backward Compatible**: No changes to provider configuration or usage
- **Future Proof**: Eliminates entire class of concurrency bugs

## Files Changed
- `internal/provider/user_resource.go`: Main fix implementation
- `internal/provider/provider.go`: Pass options instead of client
- `internal/provider/user_resource_test.go`: Fix test parameter

## Testing
The fix compiles successfully and passes existing tests. The concurrent map write errors should no longer occur during `terraform plan` operations.