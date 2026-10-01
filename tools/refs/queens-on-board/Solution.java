// Ported from the owner's backtarcking/Queens.java (nQueens and isSafe). Holes
// are skipped. Fixed: nQueens never called isSafe, so it counted every way to
// put one queen in each row (n to the power n) instead of the safe ones.
class Solution {
    private int count = 0;

    public int countPlacements(char[][] board) {
        nQueens(board, 0);
        return count;
    }

    private boolean isSafe(char[][] board, int row, int col) {
        for (int i = row - 1; i >= 0; i--) {
            if (board[i][col] == 'Q') {
                return false;
            }
        }
        for (int i = row - 1, j = col - 1; i >= 0 && j >= 0; i--, j--) {
            if (board[i][j] == 'Q') {
                return false;
            }
        }
        for (int i = row - 1, j = col + 1; i >= 0 && j < board.length; i--, j++) {
            if (board[i][j] == 'Q') {
                return false;
            }
        }
        return true;
    }

    private void nQueens(char[][] board, int row) {
        if (row == board.length) {
            count++;
            return;
        }
        for (int j = 0; j < board.length; j++) {
            if (board[row][j] == '#' || !isSafe(board, row, j)) {
                continue;
            }
            board[row][j] = 'Q';
            nQueens(board, row + 1);
            board[row][j] = '.';
        }
    }
}
