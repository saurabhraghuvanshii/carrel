#include <algorithm>
using namespace std;

class TokenBucket {
    long long capacity, refill, tokens, last = 0;

public:
    TokenBucket(int capacity, int refill) : capacity(capacity), refill(refill), tokens(capacity) {}

    bool take(int time, int n) {
        tokens = min(capacity, tokens + time / refill - last / refill);
        last = time;
        if (tokens < n) return false;
        tokens -= n;
        return true;
    }
};
