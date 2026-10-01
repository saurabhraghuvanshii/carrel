import java.util.*;

class Solution {
    public List<List<String>> groupByLetters(String[] words) {
        Map<String, List<String>> groups = new HashMap<>();
        for (String w : words) {
            int[] count = new int[26];
            for (char ch : w.toCharArray()) {
                count[ch - 'a']++;
            }
            groups.computeIfAbsent(Arrays.toString(count), k -> new ArrayList<>()).add(w);
        }
        return new ArrayList<>(groups.values());
    }
}
