#include <vector>
using namespace std;

int kthSmallest(TreeNode* root, int k) {
    vector<TreeNode*> stack;
    TreeNode* n = root;
    while (true) {
        while (n) {
            stack.push_back(n);
            n = n->left;
        }
        n = stack.back();
        stack.pop_back();
        if (--k == 0) return n->val;
        n = n->right;
    }
}
