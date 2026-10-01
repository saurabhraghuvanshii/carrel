// Ported from the owner's tries/LongestWord.java (insert, longestWord). Fixed:
// the root and the answer were static, so later calls still saw earlier
// words and answers; both now belong to each call.
class Solution {
    static class Node {
        Node[] children = new Node[26];
        boolean eow;
    }

    private Node root;
    private String ans;

    public String longestBuildable(String[] words) {
        root = new Node();
        ans = "";
        for (int i = 0; i < words.length; i++) {
            insert(words[i]);
        }
        longestWord(root, new StringBuilder(""));
        return ans;
    }

    private void insert(String word) {
        Node curr = root;
        for (int level = 0; level < word.length(); level++) {
            int idx = word.charAt(level) - 'a';
            if (curr.children[idx] == null) {
                curr.children[idx] = new Node();
            }
            curr = curr.children[idx];
        }
        curr.eow = true;
    }

    private void longestWord(Node root, StringBuilder temp) {
        if (root == null) {
            return;
        }
        for (int i = 0; i < 26; i++) {
            if (root.children[i] != null && root.children[i].eow == true) {
                char ch = (char) (i + 'a');
                temp.append(ch);
                if (temp.length() > ans.length()) {
                    ans = temp.toString();
                }
                longestWord(root.children[i], temp);
                temp.deleteCharAt(temp.length() - 1);
            }
        }
    }
}
