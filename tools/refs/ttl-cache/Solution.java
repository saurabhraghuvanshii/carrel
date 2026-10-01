import java.util.*;

class TTLCache {
    private final Map<Integer, int[]> live = new HashMap<>(); // key -> {value, expiry}
    private final PriorityQueue<int[]> byExpiry = new PriorityQueue<>((a, b) -> Integer.compare(a[0], b[0]));

    public void put(int key, int value, int ttl, int time) {
        live.put(key, new int[] {value, time + ttl});
        byExpiry.add(new int[] {time + ttl, key});
    }

    public int get(int key, int time) {
        int[] e = live.get(key);
        return e != null && e[1] > time ? e[0] : -1;
    }

    public int count(int time) {
        while (!byExpiry.isEmpty() && byExpiry.peek()[0] <= time) {
            int[] old = byExpiry.poll();
            int[] e = live.get(old[1]);
            if (e != null && e[1] == old[0]) {
                live.remove(old[1]);
            }
        }
        return live.size();
    }
}
