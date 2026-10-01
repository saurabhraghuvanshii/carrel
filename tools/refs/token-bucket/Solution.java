class TokenBucket {
    private final int capacity, refill;
    private long tokens;
    private int last;

    public TokenBucket(int capacity, int refill) {
        this.capacity = capacity;
        this.refill = refill;
        tokens = capacity;
    }

    public boolean take(int time, int n) {
        tokens = Math.min(capacity, tokens + time / refill - last / refill);
        last = time;
        if (tokens < n) {
            return false;
        }
        tokens -= n;
        return true;
    }
}
