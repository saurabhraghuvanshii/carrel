import java.util.*;

class SeatHolder {
    private final boolean[] taken; // held or sold
    private final int ttl;
    private int free;
    private final Map<Integer, int[]> holds = new HashMap<>(); // user -> {first, count, until}
    private final ArrayDeque<int[]> byExpiry = new ArrayDeque<>(); // {user, until}

    public SeatHolder(int seats, int ttl) {
        taken = new boolean[seats];
        this.ttl = ttl;
        free = seats;
    }

    private void expire(int time) {
        while (!byExpiry.isEmpty() && byExpiry.peekFirst()[1] <= time) {
            int[] old = byExpiry.pollFirst();
            int[] h = holds.get(old[0]);
            if (h != null && h[2] == old[1]) {
                giveBack(old[0]);
            }
        }
    }

    private void giveBack(int user) {
        int[] h = holds.remove(user);
        Arrays.fill(taken, h[0], h[0] + h[1], false);
        free += h[1];
    }

    public int hold(int user, int count, int time) {
        expire(time);
        if (holds.containsKey(user)) {
            return -1;
        }
        int run = 0;
        for (int s = 0; s < taken.length; s++) {
            run = taken[s] ? 0 : run + 1;
            if (run == count) {
                int first = s - count + 1;
                Arrays.fill(taken, first, s + 1, true);
                free -= count;
                holds.put(user, new int[] {first, count, time + ttl});
                byExpiry.addLast(new int[] {user, time + ttl});
                return first;
            }
        }
        return -1;
    }

    public boolean buy(int user, int time) {
        expire(time);
        return holds.remove(user) != null;
    }

    public boolean release(int user, int time) {
        expire(time);
        if (!holds.containsKey(user)) {
            return false;
        }
        giveBack(user);
        return true;
    }

    public int freeSeats(int time) {
        expire(time);
        return free;
    }
}
