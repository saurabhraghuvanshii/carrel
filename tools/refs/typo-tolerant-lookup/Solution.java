import java.util.*;

class Solution {
    public int[] nearMatches(String[] words, String[] queries) {
        Set<String> exact = new HashSet<>(Arrays.asList(words));
        Map<String, Integer> masked = new HashMap<>();
        for (String w : words) {
            for (int i = 0; i < w.length(); i++) {
                masked.merge(mask(w, i), 1, Integer::sum);
            }
        }
        int[] out = new int[queries.length];
        for (int k = 0; k < queries.length; k++) {
            String q = queries[k];
            for (int i = 0; i < q.length(); i++) {
                out[k] += masked.getOrDefault(mask(q, i), 0);
            }
            if (exact.contains(q)) {
                out[k] -= q.length() - 1;
            }
        }
        return out;
    }

    private String mask(String w, int i) {
        return w.substring(0, i) + '*' + w.substring(i + 1);
    }
}
