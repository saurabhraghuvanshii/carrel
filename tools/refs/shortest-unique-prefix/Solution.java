import java.util.*;

// Ported from the owner's tries/PrefixPrblm.java. findPrefix collects the
// prefixes instead of printing them (it finds them in alphabetical order),
// and each word then takes the one it starts with. Fixed: the root was
// static, so words from earlier calls stayed in the tree and their counts made
// prefixes too long; each call now starts a fresh tree.
class Solution {
    static class Node {
        Node[] children = new Node[26];
        boolean eow = false;
        int freq;

        public Node() {
            freq = 1;
        }
    }

    private Node root;
    private List<String> found = new ArrayList<>();

    public List<String> shortestPrefixes(String[] arr) {
        root = new Node();
        for (int i = 0; i < arr.length; i++) {
            insert(arr[i]);
        }
        root.freq = -1;
        findPrefix(root, "");
        List<String> out = new ArrayList<>();
        for (String w : arr) {
            for (String p : found) {
                if (w.startsWith(p)) {
                    out.add(p);
                    break;
                }
            }
        }
        return out;
    }

    private void insert(String word) {
        Node curr = root;
        for (int i = 0; i < word.length(); i++) {
            int idx = word.charAt(i) - 'a';
            if (curr.children[idx] == null) {
                curr.children[idx] = new Node();
            } else {
                curr.children[idx].freq++;
            }
            curr = curr.children[idx];
        }
        curr.eow = true;
    }

    private void findPrefix(Node root, String ans) {
        if (root == null) {
            return;
        }
        if (root.freq == 1) {
            found.add(ans);
            return;
        }
        for (int i = 0; i < root.children.length; i++) {
            if (root.children[i] != null) {
                findPrefix(root.children[i], ans + (char) (i + 'a'));
            }
        }
    }
}
