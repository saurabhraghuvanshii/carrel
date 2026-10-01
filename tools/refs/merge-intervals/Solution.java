import java.util.*;

class Solution {
    public int[][] merge(int[][] intervals) {
        int[][] s = intervals.clone();
        Arrays.sort(s, (a, b) -> Integer.compare(a[0], b[0]));
        List<int[]> out = new ArrayList<>();
        for (int[] iv : s) {
            if (!out.isEmpty() && iv[0] <= out.get(out.size() - 1)[1]) {
                int[] last = out.get(out.size() - 1);
                last[1] = Math.max(last[1], iv[1]);
            } else {
                out.add(new int[] {iv[0], iv[1]});
            }
        }
        return out.toArray(new int[0][]);
    }
}
