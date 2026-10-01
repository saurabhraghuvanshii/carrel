import java.util.*;

class Solution {
    public int[] topPages(int[] views, int k) {
        Map<Integer, Integer> count = new HashMap<>();
        for (int v : views) {
            count.merge(v, 1, Integer::sum);
        }
        List<Integer> pages = new ArrayList<>(count.keySet());
        pages.sort((a, b) -> count.get(a).equals(count.get(b)) ? Integer.compare(a, b) : Integer.compare(count.get(b), count.get(a)));
        int[] out = new int[Math.min(k, pages.size())];
        for (int i = 0; i < out.length; i++) {
            out[i] = pages.get(i);
        }
        return out;
    }
}
