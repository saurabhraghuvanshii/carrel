import java.util.*;

class Solution {
    private int[] parent;

    private int find(int x) {
        while (parent[x] != x) {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        return x;
    }

    public int[] keepers(String[][] contacts) {
        int n = contacts.length;
        parent = new int[n];
        for (int i = 0; i < n; i++) {
            parent[i] = i;
        }
        Map<String, Integer> firstWith = new HashMap<>();
        for (int i = 0; i < n; i++) {
            for (String address : contacts[i]) {
                Integer j = firstWith.putIfAbsent(address, i);
                if (j != null) {
                    // The lower root stays the root, so a root is its set's lowest number.
                    int a = find(i), b = find(j);
                    parent[Math.max(a, b)] = Math.min(a, b);
                }
            }
        }
        int[] out = new int[n];
        for (int i = 0; i < n; i++) {
            out[i] = find(i);
        }
        return out;
    }
}
