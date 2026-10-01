import java.util.*;

class Solution {
    public List<List<Integer>> combinations(int[] candidates, int target) {
        int[] c = candidates.clone();
        Arrays.sort(c);
        List<List<Integer>> out = new ArrayList<>();
        walk(c, 0, target, new ArrayList<>(), out);
        return out;
    }

    private void walk(int[] c, int from, int left, List<Integer> cur, List<List<Integer>> out) {
        if (left == 0) {
            out.add(new ArrayList<>(cur));
            return;
        }
        for (int i = from; i < c.length && c[i] <= left; i++) {
            cur.add(c[i]);
            walk(c, i, left - c[i], cur, out);
            cur.remove(cur.size() - 1);
        }
    }
}
