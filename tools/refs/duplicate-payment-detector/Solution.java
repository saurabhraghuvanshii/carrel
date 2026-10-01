import java.util.*;

class Solution {
    public List<Integer> duplicates(int[][] payments, int window) {
        Map<List<Integer>, Integer> last = new HashMap<>();
        List<Integer> out = new ArrayList<>();
        for (int i = 0; i < payments.length; i++) {
            int[] p = payments[i];
            Integer before = last.put(List.of(p[1], p[2], p[3]), p[0]);
            if (before != null && p[0] - before <= window) {
                out.add(i);
            }
        }
        return out;
    }
}
