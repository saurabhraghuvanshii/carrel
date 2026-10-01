import java.util.*;

class RateLimiter {
    private final int limit, window;
    private final Map<Integer, int[]> seen = new HashMap<>(); // user -> {window index, allowed}

    public RateLimiter(int limit, int window) {
        this.limit = limit;
        this.window = window;
    }

    public boolean allow(int user, int time) {
        int w = time / window;
        int[] s = seen.computeIfAbsent(user, u -> new int[] {w, 0});
        if (s[0] != w) {
            s[0] = w;
            s[1] = 0;
        }
        if (s[1] == limit) {
            return false;
        }
        s[1]++;
        return true;
    }
}
