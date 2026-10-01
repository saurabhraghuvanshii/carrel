import java.util.*;

class Solution {
    public int[] trending(int[][] posts, int window, int k) {
        int last = posts[posts.length - 1][0];
        Map<Integer, Integer> growth = new HashMap<>();
        for (int[] p : posts) {
            if (p[0] > last - window) {
                growth.merge(p[1], 1, Integer::sum);
            } else if (p[0] > last - 2L * window) {
                growth.merge(p[1], -1, Integer::sum);
            }
        }
        List<Integer> up = new ArrayList<>();
        for (Map.Entry<Integer, Integer> e : growth.entrySet()) {
            if (e.getValue() > 0) {
                up.add(e.getKey());
            }
        }
        up.sort((a, b) -> growth.get(a).equals(growth.get(b)) ? Integer.compare(a, b) : Integer.compare(growth.get(b), growth.get(a)));
        int[] out = new int[Math.min(k, up.size())];
        for (int i = 0; i < out.length; i++) {
            out[i] = up.get(i);
        }
        return out;
    }
}
