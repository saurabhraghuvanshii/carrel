class Solution {
    static class Node {
        Node[] next = new Node[26];
        int pass;
    }

    public int[] countMatches(String[] words, String[] prefixes) {
        Node root = new Node();
        for (String w : words) {
            Node cur = root;
            for (char ch : w.toCharArray()) {
                if (cur.next[ch - 'a'] == null) {
                    cur.next[ch - 'a'] = new Node();
                }
                cur = cur.next[ch - 'a'];
                cur.pass++;
            }
        }
        int[] out = new int[prefixes.length];
        for (int i = 0; i < prefixes.length; i++) {
            Node cur = root;
            for (char ch : prefixes[i].toCharArray()) {
                cur = cur.next[ch - 'a'];
                if (cur == null) {
                    break;
                }
            }
            out[i] = cur == null ? 0 : cur.pass;
        }
        return out;
    }
}
