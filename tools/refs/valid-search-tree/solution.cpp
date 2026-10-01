#include <climits>
#include <vector>
using namespace std;

// In-order values of a search tree rise strictly.
bool isSearchTree(TreeNode* root) {
    vector<TreeNode*> stack;
    long long last = LLONG_MIN;
    for (TreeNode* n = root; n || !stack.empty();) {
        while (n) {
            stack.push_back(n);
            n = n->left;
        }
        n = stack.back();
        stack.pop_back();
        if (n->val <= last) return false;
        last = n->val;
        n = n->right;
    }
    return true;
}
