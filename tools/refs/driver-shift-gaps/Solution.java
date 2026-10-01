import java.util.*;

class Solution {
    public int[][] understaffed(int[][] shifts, int m, int dayEnd) {
        TreeMap<Integer, Integer> change = new TreeMap<>();
        change.put(0, 0);
        change.put(dayEnd, 0);
        for (int[] s : shifts) {
            change.merge(s[0], 1, Integer::sum);
            change.merge(s[1], -1, Integer::sum);
        }
        List<int[]> out = new ArrayList<>();
        int on = 0;
        Integer t = change.firstKey();
        while (t != null && t < dayEnd) {
            on += change.get(t);
            Integer next = change.higherKey(t);
            if (on < m) {
                if (!out.isEmpty() && out.get(out.size() - 1)[1] == t) {
                    out.get(out.size() - 1)[1] = next;
                } else {
                    out.add(new int[] {t, next});
                }
            }
            t = next;
        }
        return out.toArray(new int[0][]);
    }
}
