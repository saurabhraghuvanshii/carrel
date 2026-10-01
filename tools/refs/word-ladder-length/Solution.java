import java.util.*;

class Solution {
    public int ladderLength(String start, String end, String[] words) {
        Set<String> left = new HashSet<>(Arrays.asList(words));
        if (!left.contains(end)) {
            return 0;
        }
        ArrayDeque<String> queue = new ArrayDeque<>();
        queue.add(start);
        left.remove(start);
        for (int steps = 1; !queue.isEmpty(); steps++) {
            for (int size = queue.size(); size > 0; size--) {
                char[] w = queue.poll().toCharArray();
                if (String.valueOf(w).equals(end)) {
                    return steps;
                }
                for (int i = 0; i < w.length; i++) {
                    char keep = w[i];
                    for (char ch = 'a'; ch <= 'z'; ch++) {
                        w[i] = ch;
                        String next = String.valueOf(w);
                        if (left.remove(next)) {
                            queue.add(next);
                        }
                    }
                    w[i] = keep;
                }
            }
        }
        return 0;
    }
}
