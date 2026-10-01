import java.util.*;

class VisitorCounter {
    private final int window;
    private final ArrayDeque<int[]> visits = new ArrayDeque<>(); // {user, time}
    private final Map<Integer, Integer> latest = new HashMap<>();

    public VisitorCounter(int window) {
        this.window = window;
    }

    public void visit(int user, int time) {
        visits.addLast(new int[] {user, time});
        latest.put(user, time);
    }

    public int unique(int time) {
        while (!visits.isEmpty() && visits.peekFirst()[1] <= time - window) {
            int[] old = visits.pollFirst();
            // The user may be gone already if they visited twice at this time.
            Integer last = latest.get(old[0]);
            if (last != null && last == old[1]) {
                latest.remove(old[0]);
            }
        }
        return latest.size();
    }
}
