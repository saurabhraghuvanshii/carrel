import java.util.*;

class Solution {
    public List<Boolean> approve(int[][] purchases, int[][] requests, int window) {
        Map<Integer, Integer> bought = new HashMap<>();
        for (int[] p : purchases) {
            bought.put(p[0], p[1]);
        }
        Set<Integer> refunded = new HashSet<>();
        List<Boolean> out = new ArrayList<>();
        for (int[] r : requests) {
            Integer day = bought.get(r[0]);
            boolean ok = day != null && r[1] >= day && r[1] - day <= window && !refunded.contains(r[0]);
            if (ok) {
                refunded.add(r[0]);
            }
            out.add(ok);
        }
        return out;
    }
}
