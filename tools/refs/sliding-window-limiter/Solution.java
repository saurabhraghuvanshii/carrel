import java.util.*;

class SlidingLimiter {
    private final int limit, window;
    private final Map<Integer, ArrayDeque<Integer>> recent = new HashMap<>();

    public SlidingLimiter(int limit, int window) {
        this.limit = limit;
        this.window = window;
    }

    public boolean allow(int user, int time) {
        ArrayDeque<Integer> q = recent.computeIfAbsent(user, u -> new ArrayDeque<>());
        while (!q.isEmpty() && q.peekFirst() <= time - window) {
            q.pollFirst();
        }
        if (q.size() == limit) {
            return false;
        }
        q.addLast(time);
        return true;
    }
}
