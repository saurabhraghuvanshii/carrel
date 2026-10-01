import java.util.*;

class LFUCache {
    private final int capacity;
    private final Map<Integer, int[]> entries = new HashMap<>(); // key -> {value, uses}
    private final Map<Integer, LinkedHashSet<Integer>> byUses = new HashMap<>();
    private int lowest;

    public LFUCache(int capacity) {
        this.capacity = capacity;
    }

    public int get(int key) {
        int[] e = entries.get(key);
        if (e == null) {
            return -1;
        }
        use(key, e);
        return e[0];
    }

    public void put(int key, int value) {
        int[] e = entries.get(key);
        if (e != null) {
            e[0] = value;
            use(key, e);
            return;
        }
        if (entries.size() == capacity) {
            Iterator<Integer> oldest = byUses.get(lowest).iterator();
            entries.remove(oldest.next());
            oldest.remove();
        }
        entries.put(key, new int[] {value, 1});
        byUses.computeIfAbsent(1, u -> new LinkedHashSet<>()).add(key);
        lowest = 1;
    }

    private void use(int key, int[] e) {
        LinkedHashSet<Integer> set = byUses.get(e[1]);
        set.remove(key);
        if (set.isEmpty() && lowest == e[1]) {
            lowest++;
        }
        e[1]++;
        byUses.computeIfAbsent(e[1], u -> new LinkedHashSet<>()).add(key);
    }
}
