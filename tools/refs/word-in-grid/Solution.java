class Solution {
    public boolean hasWord(char[][] grid, String word) {
        for (int r = 0; r < grid.length; r++) {
            for (int c = 0; c < grid[0].length; c++) {
                if (walk(grid, word, r, c, 0)) {
                    return true;
                }
            }
        }
        return false;
    }

    private boolean walk(char[][] g, String word, int r, int c, int k) {
        if (r < 0 || r >= g.length || c < 0 || c >= g[0].length || g[r][c] != word.charAt(k)) {
            return false;
        }
        if (k == word.length() - 1) {
            return true;
        }
        char keep = g[r][c];
        g[r][c] = '#';
        boolean found = walk(g, word, r + 1, c, k + 1) || walk(g, word, r - 1, c, k + 1)
                || walk(g, word, r, c + 1, k + 1) || walk(g, word, r, c - 1, k + 1);
        g[r][c] = keep;
        return found;
    }
}
